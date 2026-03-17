package com.videoplatform.controller;

import com.videoplatform.service.SearchService;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/search")
@RequiredArgsConstructor
public class SearchController {

    private final SearchService searchService;

    @GetMapping
    public ResponseEntity<?> search(
            @RequestParam String q,
            @RequestParam(defaultValue = "video") String type,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        
        if ("video".equals(type)) {
            Page<SearchService.VideoSearchResult> results = searchService.searchVideos(q, page, size);
            return ResponseEntity.ok(Map.of(
                "results", results.getContent(),
                "total", results.getTotalElements(),
                "page", page,
                "page_size", size
            ));
        } else if ("user".equals(type)) {
            Page<SearchService.UserSearchResult> results = searchService.searchUsers(q, page, size);
            return ResponseEntity.ok(Map.of(
                "results", results.getContent(),
                "total", results.getTotalElements(),
                "page", page,
                "page_size", size
            ));
        } else {
            Map<String, Object> results = searchService.searchAll(q, page, size);
            return ResponseEntity.ok(results);
        }
    }

    @GetMapping("/videos")
    public ResponseEntity<?> searchVideos(
            @RequestParam String q,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        Page<SearchService.VideoSearchResult> results = searchService.searchVideos(q, page, size);
        return ResponseEntity.ok(Map.of(
            "results", results.getContent(),
            "total", results.getTotalElements(),
            "page", page,
            "page_size", size
        ));
    }

    @GetMapping("/users")
    public ResponseEntity<?> searchUsers(
            @RequestParam String q,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        Page<SearchService.UserSearchResult> results = searchService.searchUsers(q, page, size);
        return ResponseEntity.ok(Map.of(
            "results", results.getContent(),
            "total", results.getTotalElements(),
            "page", page,
            "page_size", size
        ));
    }
}
