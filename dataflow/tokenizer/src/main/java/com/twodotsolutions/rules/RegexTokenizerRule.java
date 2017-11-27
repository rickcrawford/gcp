package com.twodotsolutions.rules;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;


public class RegexTokenizerRule implements Rule {
  /**
   * 
   */
  private static final long serialVersionUID = 2579015161676076528L;
  private Pattern pattern;
  private boolean includeOtherTokens;

  public RegexTokenizerRule(String regex, boolean includeOtherTokens) {
    this.pattern = Pattern.compile(regex);
    this.includeOtherTokens = includeOtherTokens;
  }

  @Override
  public List<String> apply(List<String> tokens) {
    List<String> results = new ArrayList<>(tokens.size() * 2);
    StringBuilder sb;
    Matcher m;
    String otherToken, remaining;

    for (String token : tokens) {
      sb = new StringBuilder(token);
      m = pattern.matcher(token);
      int lastEnd = 0;
      while (m.find()) {
        results.add(token.substring(m.start() + 1, m.end() - 1).trim());
        if (includeOtherTokens) {
          otherToken = token.substring(lastEnd, m.start()).trim();
          if (!otherToken.isEmpty()) {
            results.add(otherToken);
          }
        }
        if (lastEnd < sb.length()) {
          sb.delete(lastEnd, m.end() + 1);
        }
        lastEnd = m.end() + 1;
      }
      remaining = sb.toString().trim();
      if (!remaining.isEmpty() && includeOtherTokens) {
        results.add(remaining);
      }
    }

    return results;

  }
}
