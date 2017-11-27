package com.twodotsolutions.rules;

import java.util.ArrayList;
import java.util.List;


public class RegexReplaceAllRule implements Rule {
  /**
   * 
   */
  private static final long serialVersionUID = 1140375059990761655L;
  private String regex;
  private String replacement;

  public RegexReplaceAllRule(String regex, String replacement) {
    this.regex = regex;
    this.replacement = replacement;
  }

  @Override
  public List<String> apply(List<String> tokens) {
    List<String> results = new ArrayList<>(tokens.size());
    String clean;
    for (String token : tokens) {
      clean = token.replaceAll(this.regex, this.replacement).trim();
      if (!clean.isEmpty()) {
        results.add(clean);
      }
    }
    return results;
  }
}
