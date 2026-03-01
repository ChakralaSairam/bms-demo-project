package model

// Movie represents a movie in the inventory catalog.
//
// Java equivalent:
//   @Entity
//   public class Movie {
//       private String id;
//       private String title;
//       private String genre;
//       private int releaseYear;
//       private double rating;
//   }
type Movie struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Genre       string  `json:"genre"`
	ReleaseYear int     `json:"release_year"`
	Rating      float64 `json:"rating"`
}
