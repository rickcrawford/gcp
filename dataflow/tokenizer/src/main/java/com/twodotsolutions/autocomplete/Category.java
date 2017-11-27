package com.twodotsolutions.autocomplete;

import java.io.Serializable;

public class Category implements Serializable {
  /**
   * 
   */
  private static final long serialVersionUID = 909613729868529503L;
  
  private String id;
  private String name;
  
  public Category() {}

  public String getId() {
    return id;
  }

  public void setId(String id) {
    this.id = id;
  }

  public String getName() {
    return name;
  }

  public void setName(String name) {
    this.name = name;
  }
  
}
