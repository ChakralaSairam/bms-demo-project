package com.test.user.client;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.*;

@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class SeatLockResponse {
    private String status;
    private String message;
    
    @JsonProperty("show_id")
    private Long showId;
    
    @JsonProperty("seat_id")
    private String seatId;
    
    @JsonProperty("user_id")
    private Long userId;
    
    @JsonProperty("held_for")
    private String heldFor;
}
