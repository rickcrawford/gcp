package com.twodotsolutions.autocomplete;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.TimeZone;

import org.apache.beam.sdk.Pipeline;
import org.apache.beam.sdk.io.FileBasedSink;
import org.apache.beam.sdk.io.TextIO;
import org.apache.beam.sdk.io.TextIO.Read;
import org.apache.beam.sdk.metrics.Counter;
import org.apache.beam.sdk.metrics.Metrics;
import org.apache.beam.sdk.options.Default;
import org.apache.beam.sdk.options.Description;
import org.apache.beam.sdk.options.PipelineOptions;
import org.apache.beam.sdk.options.PipelineOptionsFactory;
import org.apache.beam.sdk.transforms.Count;
import org.apache.beam.sdk.transforms.DoFn;
import org.apache.beam.sdk.transforms.MapElements;
import org.apache.beam.sdk.transforms.PTransform;
import org.apache.beam.sdk.transforms.ParDo;
import org.apache.beam.sdk.transforms.SimpleFunction;
import org.apache.beam.sdk.values.KV;
import org.apache.beam.sdk.values.PCollection;
import org.joda.time.format.ISODateTimeFormat;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.twodotsolutions.rules.LowercaseRule;
import com.twodotsolutions.rules.RegexSplitRule;
import com.twodotsolutions.rules.RegexTokenizerRule;
import com.twodotsolutions.rules.Rule;

//
// ngrams
// term
// a: count(10)
//
public class DataflowJob {

  static class ProductFn extends DoFn<String, Product> {
    private static final long serialVersionUID = 4618262921094558199L;
    private final Counter emptyLines = Metrics.counter(ProductFn.class, "emptyLines");

    @ProcessElement
    public void processElement(ProcessContext c) {
      String element = c.element().trim();
      // if the file is clean we don't need to do this...
      // this is to account for the file being one large JSON blob
      if (!element.isEmpty()) {
        if (element.charAt(0) == '[') {
          element = element.substring(1, element.length());
        }
        if (element.charAt(element.length() - 1) == ']') {
          element = element.substring(0, element.length() - 1);
        }
        if (element.charAt(element.length() - 1) == ',') {
          element = element.substring(0, element.length() - 1);
        }
      }

      // ignore empties
      if (element.isEmpty()) {
        emptyLines.inc();
        return;
      }

      Gson gson = new GsonBuilder().create();
      Product p = gson.fromJson(element, Product.class);
      if (p == null) {
        emptyLines.inc();
        return;
      }

      c.output(p);
    }
  }


  static class ExtractNgramsFn extends DoFn<Product, KV<String, String>> {
    private static final long serialVersionUID = 810714365050459782L;
    private final Counter matchedWords = Metrics.counter(ExtractNgramsFn.class, "matchedWords");

    private List<Rule> rules;
    private int maxNgramSize;
    private int minWordLength;
    private boolean includeCategories;
    private String splitRegex;

    public ExtractNgramsFn(String splitRegex, List<Rule> rules, int maxNgramSize, int minWordLength,
        boolean includeCategories) {
      this.rules = rules;
      this.maxNgramSize = maxNgramSize;
      this.minWordLength = minWordLength;
      this.includeCategories = includeCategories;
      this.splitRegex = splitRegex;
    }

    @ProcessElement
    public void processElement(ProcessContext c) {
      Product p = c.element();
      List<String> tokens = new ArrayList<>();
      tokens.add(p.getName());
      for (Rule rule : rules) {
        tokens = rule.apply(tokens);
      }
      StringBuilder sb = new StringBuilder();
      List<String> ngrams;
      String[] words;

      for (String token : tokens) {
        ngrams = new ArrayList<>();
        words = token.split(this.splitRegex);
        for (int i = 0; i < words.length; i++) {
          String word1 = words[i].trim();
          if (word1.isEmpty()) {
            continue;
          }

          // avoid single characters...
          if (word1.length() >= this.minWordLength) {
            ngrams.add(word1);
          }

          sb.setLength(0);
          sb.append(word1);

          int endIndex = i + this.maxNgramSize;
          if (endIndex > words.length) {
            endIndex = words.length;
          }

          for (int j = i + 1; j < endIndex; j++) {
            String word2 = words[j].trim();
            if (word2.isEmpty()) {
              continue;
            }
            sb.append(" ");
            sb.append(word2.trim());
            ngrams.add(sb.toString());
          }
        }

        Category firstCategory = null;
        if (includeCategories) {
          firstCategory = p.getCategory().iterator().next();
        }

        for (String ngram : ngrams) {
          ngram = ngram.trim();
          if (!ngram.isEmpty()) {
            c.output(KV.of(ngram, ""));
            if (firstCategory != null) {
              c.output(KV.of(ngram, firstCategory.getName()));
            }
            matchedWords.inc();
          }
        }
      }
    }
  }


  public static class FormatAsTextFn extends SimpleFunction<KV<KV<String, String>, Long>, String> {
    private static final long serialVersionUID = -8437112932968569136L;

    @Override
    public String apply(KV<KV<String, String>, Long> input) {
      KV<String, String> key = input.getKey();
      Long value = input.getValue();

      Gson gson = new GsonBuilder().create();
      Result result = new Result();
      result.setCategory(key.getValue());
      result.setKeyword(key.getKey());
      result.setCount(value.intValue());

      return gson.toJson(result);
    }
  }


  static class FilterSingleCountFn<T> extends DoFn<KV<T, Long>, KV<T, Long>> {

    private static final long serialVersionUID = -3557953772166303839L;

    private int size;

    public FilterSingleCountFn(int size) {
      this.size = size;
    }

    @ProcessElement
    public void processElement(ProcessContext c) {
      KV<T, Long> row = c.element();
      if (row.getValue() >= this.size) {
        c.output(row);
      }
    }
  }


  static class CreateTokens
      extends PTransform<PCollection<String>, PCollection<KV<KV<String, String>, Long>>> {


    /**
     * Worth looking at:
     * https://github.com/apache/beam/blob/master/examples/java/src/main/java/org/apache/beam/examples/complete/TfIdf.java
     */
    private static final long serialVersionUID = -3270608019547038296L;

    @Override
    public PCollection<KV<KV<String, String>, Long>> expand(PCollection<String> lines) {

      // --> line --> product
      PCollection<Product> products = lines.apply("LoadProducts", ParDo.of(new ProductFn()));

      List<Rule> rules = new ArrayList<>();
      rules.add(new RegexSplitRule("\\s+-\\s+"));
      rules.add(new RegexTokenizerRule("\"(\\w(.*?)?)\"", true)); // quoted strings
      rules.add(new RegexTokenizerRule("\\(([^\\)]+)\\)", true)); // parens
      rules.add(new RegexTokenizerRule("\\[([^\\]]+)\\]", true)); // brackets
      rules.add(new LowercaseRule());

      // product --> {term: category}
      PCollection<KV<String, String>> terms = products.apply("CreateTerms",
          ParDo.of(new ExtractNgramsFn("[^\\p{L}\\p{N}\\/\\.]+", rules, 4, 2, true)));

      // {term: category} --> {{term: category}, count}
      PCollection<KV<KV<String, String>, Long>> termCategoryCounts =
          terms.apply("CountTermCategory", Count.<KV<String, String>>perElement());


      // Filter
      PCollection<KV<KV<String, String>, Long>> filterCounts = termCategoryCounts
          .apply("FilterCounts", ParDo.of(new FilterSingleCountFn<KV<String, String>>(2)));


      return filterCounts;

      // return products;
    }
  }


  public interface JobOptions extends PipelineOptions {

    /**
     * By default, this example reads from a public dataset containing the text of King Lear. Set
     * this option to choose a different input file or glob.
     */
    // @Default.String("gs://dataprep-staging-3600cb4e-4a05-47b2-b3dc-4b7b6aaa3a67/rick.crawford@gmail.com/jobrun/products_clean.json")

    @Description("Path of the file to read from")
    @Default.String("gs://typeahead-catalogs/bestbuy/products.json.gz")
    String getInputFile();

    void setInputFile(String value);

    /**
     * Set this required option to specify where to write the output.
     */
    @Description("Path of the file to write to")
    @Default.String("gs://typeahead-catalogs/output/")
    String getOutputDirectory();

    void setOutputDirectory(String value);
  }


  // Run the application
  public static void main(String[] args) {

    JobOptions options =
        PipelineOptionsFactory.fromArgs(args).withValidation().as(JobOptions.class);

    Read read = TextIO.read().from(options.getInputFile());
    if (options.getInputFile().endsWith(".gz")) {
      read = read.withCompressionType(TextIO.CompressionType.GZIP);
    }
    
    String output = options.getOutputDirectory();
    if (!output.endsWith("/")) {
      output += "/";
    }

    TimeZone tz = TimeZone.getTimeZone("UTC");
    DateFormat df = new SimpleDateFormat("yyyyMMddHHmm");
    df.setTimeZone(tz);
    
    output = output + df.format(new Date()) + "/keywords";

    Pipeline p = Pipeline.create(options);
    p.apply("ReadLines", read)

        .apply(new CreateTokens()).apply(MapElements.via(new FormatAsTextFn()))
        .apply("WriteCounts", TextIO.write().to(output)
            .withWritableByteChannelFactory(FileBasedSink.CompressionType.GZIP));

    p.run().waitUntilFinish();



  }
}
