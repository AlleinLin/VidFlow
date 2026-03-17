package com.videoplatform.service;

import com.videoplatform.domain.entity.Comment;
import com.videoplatform.domain.entity.Like;
import com.videoplatform.domain.entity.Video;
import com.videoplatform.domain.repository.CommentRepository;
import com.videoplatform.domain.repository.LikeRepository;
import com.videoplatform.domain.repository.VideoRepository;
import com.videoplatform.dto.comment.CommentListResponse;
import com.videoplatform.dto.comment.CommentRequest;
import com.videoplatform.dto.comment.CommentResponse;
import com.videoplatform.exception.BusinessException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class CommentService {

    private final CommentRepository commentRepository;
    private final VideoRepository videoRepository;
    private final LikeRepository likeRepository;

    @Transactional
    public CommentResponse create(Long userId, CommentRequest request) {
        Video video = videoRepository.findById(request.getVideoId())
                .orElseThrow(() -> new BusinessException("Video not found"));

        Comment.CommentBuilder builder = Comment.builder()
                .video(video)
                .user(com.videoplatform.domain.entity.User.builder().id(userId).build())
                .content(request.getContent());

        if (request.getParentId() != null) {
            Comment parent = commentRepository.findById(request.getParentId())
                    .orElseThrow(() -> new BusinessException("Parent comment not found"));

            builder.parent(parent);
            if (parent.getParent() != null) {
                builder.root(parent.getRoot());
            } else {
                builder.root(parent);
            }
        }

        Comment comment = commentRepository.save(builder.build());
        videoRepository.incrementCommentCount(request.getVideoId());

        log.info("Comment created: commentId={}, videoId={}", comment.getId(), request.getVideoId());

        return mapToResponse(comment);
    }

    @Transactional(readOnly = true)
    public CommentListResponse getByVideoId(Long videoId, int page, int size) {
        PageRequest pageable = PageRequest.of(page - 1, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        Page<Comment> commentPage = commentRepository.findByVideoIdAndStatusAndParentIsNull(
                videoId, Comment.CommentStatus.VISIBLE, pageable);

        List<CommentResponse> comments = commentPage.getContent().stream()
                .map(this::mapToResponse)
                .toList();

        return CommentListResponse.builder()
                .comments(comments)
                .total(commentPage.getTotalElements())
                .build();
    }

    @Transactional
    public void delete(Long id, Long userId) {
        Comment comment = commentRepository.findById(id)
                .orElseThrow(() -> new BusinessException("Comment not found"));

        if (!comment.getUser().getId().equals(userId)) {
            throw new BusinessException("Unauthorized");
        }

        comment.setStatus(Comment.CommentStatus.DELETED);
        commentRepository.save(comment);

        videoRepository.decrementCommentCount(comment.getVideo().getId());

        log.info("Comment deleted: commentId={}", id);
    }

    @Transactional
    public void like(Long userId, Long commentId) {
        if (likeRepository.existsByUserIdAndTargetIdAndType(userId, commentId, Like.LikeType.COMMENT)) {
            throw new BusinessException("Already liked");
        }

        Like like = Like.builder()
                .user(com.videoplatform.domain.entity.User.builder().id(userId).build())
                .targetId(commentId)
                .type(Like.LikeType.COMMENT)
                .build();
        likeRepository.save(like);

        commentRepository.incrementLikeCount(commentId);
    }

    @Transactional
    public void unlike(Long userId, Long commentId) {
        if (!likeRepository.existsByUserIdAndTargetIdAndType(userId, commentId, Like.LikeType.COMMENT)) {
            throw new BusinessException("Not liked");
        }

        likeRepository.deleteByUserIdAndTargetIdAndType(userId, commentId, Like.LikeType.COMMENT);
        commentRepository.decrementLikeCount(commentId);
    }

    private CommentResponse mapToResponse(Comment comment) {
        return CommentResponse.builder()
                .id(comment.getId())
                .videoId(comment.getVideo().getId())
                .userId(comment.getUser().getId())
                .parentId(comment.getParent() != null ? comment.getParent().getId() : null)
                .rootId(comment.getRoot() != null ? comment.getRoot().getId() : null)
                .content(comment.getContent())
                .likeCount(comment.getLikeCount())
                .status(comment.getStatus().name())
                .createdAt(comment.getCreatedAt())
                .build();
    }
}
