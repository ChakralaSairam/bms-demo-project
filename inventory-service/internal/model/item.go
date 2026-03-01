package model

// In Java, this would be a class:
//
//   public class Item {
//       private String id;
//       private String name;
//       private int quantity;
//       private float64 price;
//       // getters, setters, constructors...
//   }
//
// In Go, we use a "struct" instead of a class.
// There are NO getters/setters by default — fields are accessed directly.
// Exported (public) fields start with an Uppercase letter.
// Unexported (private) fields start with a lowercase letter.

// Item represents a product in the inventory.
type Item struct {
	ID       string  `json:"id"`       // exported = public in Java
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// The backtick annotations (e.g., `json:"id"`) are called "struct tags".
// They are similar to Java annotations like @JsonProperty("id") in Jackson.
