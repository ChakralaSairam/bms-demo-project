package com.test.user.client;

import com.test.user.config.FeignConfig;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.*;

@FeignClient(name = "inventory-service", url = "${inventory.service.url}", configuration = FeignConfig.class)
public interface InventoryClient {
    
    @PostMapping("/inventory/lock")
    SeatLockResponse lockSeat(@RequestHeader("Authorization") String token, @RequestBody SeatLockRequest request);
    
    @DeleteMapping("/inventory/lock")
    SeatUnlockResponse unlockSeat(
        @RequestHeader("Authorization") String token,
        @RequestParam("show_id") Long showId,
        @RequestParam("seat_id") String seatId
    );
}
