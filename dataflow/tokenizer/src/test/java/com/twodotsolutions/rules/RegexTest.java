package com.twodotsolutions.rules;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.apache.beam.sdk.transforms.Regex;
import org.junit.Test;

public class RegexTest {

  static String[] testTitles = new String[] {
      "product - avr-x4200w 1645w 7.2-ch. 4k ultra hd and 3d pass-through a/v home theater receiver - feature",
      "product - sleeve for surface pro 3/pro 4 and most 12\" tablets",
      "product - carbon 9.8' usb a-to-apple® \"ipod®/ipad®\" cable - feature",
      "product - fujinon xf \"14mm f/2.8\" ultrawide-angle lens for fujifilm \"x-mount system\" cameras - feature",
      "product - hard floor wipes for dyson hard dc56 vacuums (1 pack of 12 wipes) - feature",
      "product - ps 6\" x 9\" coaxial speakers (pair) - feature",
      "product - [ovrmld] case for apple® iphone® se, 5s and 5 - feature"};



  @Test
  public void test() {
    List<Rule> rules = new ArrayList<>();
//    rules.add(new RegexSplitRule("\\s+-\\s+"));
//    rules.add(new RegexTokenizerRule("\"(\\w(.*?)?)\"", true)); // quoted strings
    rules.add(new RegexTokenizerRule("\\(([^\\)]+)\\)", true)); // parens
//    rules.add(new RegexTokenizerRule("\\[([^\\]]+)\\]", true)); // brackets
    rules.add(new LowercaseRule());
    
    for (String s : testTitles) {
      List<String> results = new ArrayList<>();
      results.add(s);
      for (Rule r : rules) {
        results = r.apply(results);
      }
      System.out.println(s + "\n----\n" + results + "\n\n");
    }
    
    
  }

}
