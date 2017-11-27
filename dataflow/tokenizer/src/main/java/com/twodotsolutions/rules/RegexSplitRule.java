package com.twodotsolutions.rules;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.regex.Pattern;


public class RegexSplitRule implements Rule {
  /**
   * 
   */
  private static final long serialVersionUID = 7001171150811773267L;
  private Pattern pattern;

  public RegexSplitRule(String regex) {
    this.pattern = Pattern.compile(regex);
  }

  @Override
  public List<String> apply(List<String> tokens) {
    List<String> results = new ArrayList<>();
    for (String token : tokens) {
      if (token != null && !token.isEmpty()) {
        results.addAll(Arrays.asList(this.pattern.split(token)));
      }
    }
    return results;
  }
}
