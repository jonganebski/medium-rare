export const filterBlack = document.getElementById("filter-black");
export const editorReadOnlyHeader = document.getElementById(
  "editor-readOnly__header"
);
export const fixedAuthorInfo = document.getElementById("fixed-authorInfo");
export const likedContainer = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__liked-container"
);
export const seeCommentDiv = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__comment-container"
);
export const bookmarkContainer = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__bookmark-container"
);
export const followBtn = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__follow-btn"
) as HTMLButtonElement | null;
export const followingBtn = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__following-btn"
) as HTMLButtonElement | null;

export const commentCountDisplay = seeCommentDiv?.querySelector("span");
export const commentDrawer = document.getElementById("drawer-comment");
export const commentDrawerCloseIcon = commentDrawer?.querySelector(
  ".drawer-comment__close-icon"
);
export const preparedCommentBox = commentDrawer?.querySelector(
  ".add-comment__text"
) as HTMLParagraphElement | null;

export const commentsUlEl = commentDrawer?.querySelector("ul");
