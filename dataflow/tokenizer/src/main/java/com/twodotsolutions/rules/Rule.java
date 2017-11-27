package com.twodotsolutions.rules;

import java.io.Serializable;
import java.util.List;

public interface Rule extends Serializable {
  List<String> apply(List<String> tokens);
}
