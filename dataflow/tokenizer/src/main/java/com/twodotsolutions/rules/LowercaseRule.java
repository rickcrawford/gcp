package com.twodotsolutions.rules;

import java.util.ArrayList;
import java.util.List;

public class LowercaseRule implements Rule {

  /**
   * 
   */
  private static final long serialVersionUID = -5771385948772142667L;

  @Override
  public List<String> apply(List<String> tokens) {
    List<String> results = new ArrayList<>(tokens.size());
    for (String token : tokens) {
      token = token.toLowerCase().trim();
      if (!token.isEmpty()) {
        results.add(token);
      }
    }
    return results;
  }

}
