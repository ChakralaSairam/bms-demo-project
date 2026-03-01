package com.test.user.repository;

import com.test.user.models.Booking;
import com.test.user.models.BookingStatus;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface BookingRepository extends JpaRepository<Booking, UUID> {
    
    @Query("SELECT b.seatId FROM Booking b WHERE b.showId = :showId AND b.status = :status")
    List<String> findSeatIdsByShowIdAndStatus(@Param("showId") Long showId, @Param("status") BookingStatus status);
}
