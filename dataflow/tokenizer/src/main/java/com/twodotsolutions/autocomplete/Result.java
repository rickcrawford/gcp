package com.twodotsolutions.autocomplete;

import java.io.Serializable;

public class Result implements Serializable {
  private static final long serialVersionUID = -7500512035191225269L;
  private String keyword;
  private String category;
  private int count;
  
  public String getKeyword() {
    return keyword;
  }
  public void setKeyword(String keyword) {
    this.keyword = keyword;
  }
  public String getCategory() {
    return category;
  }
  public void setCategory(String category) {
    this.category = category;
  }
  public int getCount() {
    return count;
  }
  public void setCount(int count) {
    this.count = count;
  }
  
  

}
