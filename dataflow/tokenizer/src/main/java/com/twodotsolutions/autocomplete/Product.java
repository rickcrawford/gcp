package com.twodotsolutions.autocomplete;

import java.io.Serializable;
import java.util.List;

public class Product implements Serializable {
  /**
   * 
   */
  private static final long serialVersionUID = -7634935292796246095L;
  
  private String sku;
  private String name;
  private String description;
  private String price;
  private String shipping;
  private String upc;
  private String manufacturer;
  private String model;
  private String url;
  private String image;
  private List<Category> category;

  public String getSku() {
    return sku;
  }

  public void setSku(String sku) {
    this.sku = sku;
  }

  public String getName() {
    return name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getDescription() {
    return description;
  }

  public void setDescription(String description) {
    this.description = description;
  }

  public String getPrice() {
    return price;
  }

  public void setPrice(String price) {
    this.price = price;
  }

  public String getShipping() {
    return shipping;
  }

  public void setShipping(String shipping) {
    this.shipping = shipping;
  }

  public String getUpc() {
    return upc;
  }

  public void setUpc(String upc) {
    this.upc = upc;
  }

  public String getManufacturer() {
    return manufacturer;
  }

  public void setManufacturer(String manufacturer) {
    this.manufacturer = manufacturer;
  }

  public String getModel() {
    return model;
  }

  public void setModel(String model) {
    this.model = model;
  }

  public String getUrl() {
    return url;
  }

  public void setUrl(String url) {
    this.url = url;
  }

  public String getImage() {
    return image;
  }

  public void setImage(String image) {
    this.image = image;
  }

  public List<Category> getCategory() {
    return category;
  }

  public void setCategory(List<Category> category) {
    this.category = category;
  }

  public Product() {}

  @Override
  public String toString() {
    return sku + " - " + name;
  }
}
