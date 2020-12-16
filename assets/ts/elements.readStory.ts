export const filterBlack = document.getElementById("filter-black");

export const editorReadOnlyHeader = document.getElementById(
  "editor-readOnly__header"
);

const editorReadOnlyAuthorArea = document.getElementById(
  "editor-readOnly__authorArea"
);
export const readTimeSpan = editorReadOnlyAuthorArea?.querySelector(
  ".editor-readOnly__authorArea__readTime"
) as HTMLSpanElement | null;
export const deleteStoryBtn = editorReadOnlyAuthorArea?.querySelector(
  ".editor-readOnly__authorArea__delete"
);
export const pickStoryBtn = editorReadOnlyAuthorArea?.querySelector(
  ".editor-readOnly__pick"
);
export const unpickStoryBtn = editorReadOnlyAuthorArea?.querySelector(
  ".editor-readOnly__unpick"
);

export const editorReadOnlyBody = document.getElementById(
  "editor-readOnly__body"
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
export const followersCountDisplay = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__followerCount"
) as HTMLElement | null;
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
export const commentDrawerCommentCount = commentDrawer?.querySelector(
  ".drawer-comment__commentCount"
) as HTMLSpanElement | null;
export const preparedCommentBox = commentDrawer?.querySelector(
  ".add-comment__text"
) as HTMLParagraphElement | null;
export const cancelCommentBtn = commentDrawer?.querySelector(
  ".add-comment__cancel-btn"
);
export const addCommentBtn = commentDrawer?.querySelector(
  ".add-comment__add-btn"
);
export const commentsUlEl = commentDrawer?.querySelector("ul");
