package com.test.user.models;

import jakarta.persistence.*;
import lombok.*;
import java.math.BigDecimal;
import java.time.LocalDateTime;

@Entity
@Table(name = "shows", indexes = {
    @Index(name = "idx_movie_id", columnList = "movieId"),
    @Index(name = "idx_theater_id", columnList = "theaterId")
})
@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Show {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private Long movieId;

    @Column(nullable = false)
    private Long theaterId;

    @Column(nullable = false)
    private LocalDateTime startTime;

    @Column(nullable = false)
    private BigDecimal price;
}
