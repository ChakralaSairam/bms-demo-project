package com.test.user.models;

import jakarta.persistence.*;
import lombok.*;

@Entity
@Table(name = "theaters")
@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Theater {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String name;

    private String city;
}
