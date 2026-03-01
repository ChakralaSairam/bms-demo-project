package com.test.user.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.*;

@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class MovieDTO {
    private Long id;
    
    @NotBlank(message = "Title is required")
    private String title;
    
    private String genre;
    
    @NotNull(message = "Duration is required")
    private Integer duration;
    
    private String language;
}
