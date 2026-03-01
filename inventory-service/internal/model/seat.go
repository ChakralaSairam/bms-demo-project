package model

// SeatLockRequest is the JSON body sent by the client when locking a seat.
//
// Java equivalent (with Lombok):
//   @Data
//   public class SeatLockRequest {
//       @JsonProperty("show_id")  private Long showId;
//       @JsonProperty("seat_id")  private String seatId;
//       @JsonProperty("user_id")  private Long userId;
//   }
//
// The `json:"..."` struct tags tell Go's JSON decoder what field name to
// expect in the request body — same as @JsonProperty in Jackson.
type SeatLockRequest struct {
	ShowID int64  `json:"show_id"`
	SeatID string `json:"seat_id"`
	UserID int64  `json:"user_id"`
}
