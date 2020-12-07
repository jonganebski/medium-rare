import EditorJS from "@editorjs/editorjs";
import Header from "@editorjs/header";
import CodeTool from "@editorjs/code";
import ImageTool from "@editorjs/image";
import { BASE_URL } from "./constants";
import Axios from "axios";

const filterBlack = document.getElementById("filter-black");
const editorReadOnlyHeader = document.getElementById("editor-readOnly__header");
const fixedAuthorInfo = document.getElementById("fixed-authorInfo");
const likedContainer = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__liked-container"
);
const seeCommentDiv = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__comment-container"
);
const bookmarkContainer = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__bookmark-container"
);
const followBtn = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__follow-btn"
) as HTMLButtonElement | null;
const followingBtn = fixedAuthorInfo?.querySelector(
  ".fixed-authorInfo__following-btn"
) as HTMLButtonElement | null;

export const commentCountDisplay = seeCommentDiv?.querySelector("span");
export const commentDrawer = document.getElementById("drawer-comment");
const commentDrawerCloseIcon = commentDrawer?.querySelector(
  ".drawer-comment__close-icon"
);
export const preparedCommentBox = commentDrawer?.querySelector(
  ".add-comment__text"
) as HTMLParagraphElement | null;

const commentsUlEl = commentDrawer?.querySelector("ul");

export const drawNewComment = (comment: any) => {
  if (commentsUlEl) {
    const {
      commentId,
      avatarUrl,
      createdAt,
      text,
      username,
      isAuthorized,
    } = comment;
    const textEl = document.createElement("p");
    const timestampEl = document.createElement("div");
    const usernameEl = document.createElement("div");
    const infoEl = document.createElement("div");
    const avatarImgEl = document.createElement("img");
    const avatarFrameEl = document.createElement("div");
    let deleteIcon;
    const headerEl = document.createElement("header");
    const liEl = document.createElement("li");
    textEl.className = "comment__text";
    timestampEl.className = "comment__timestamp";
    usernameEl.className = "comment__creator-name";
    infoEl.className = "comment__info _flex-c-sb";
    avatarImgEl.className = "_avatar-img";
    avatarFrameEl.className = "_avatar-frame";
    headerEl.className = "comment__header _flex-cs";
    liEl.className = "drawer-comment__comment-container";
    liEl.id = commentId;
    textEl.innerText = text;
    timestampEl.innerText = createdAt;
    usernameEl.innerText = username;
    avatarImgEl.src = avatarUrl;
    infoEl.append(usernameEl);
    infoEl.append(timestampEl);
    avatarFrameEl.append(avatarImgEl);
    if (isAuthorized) {
      deleteIcon = document.createElement("i");
      deleteIcon.className = "far fa-trash-alt comment__delete-icon";
      deleteIcon.addEventListener("click", async () => {
        const { status } = await Axios.delete(
          BASE_URL + `/api/comment/${commentId}`
        );
        console.log(status);
        if (status === 200) {
          liEl.remove();
          if (commentCountDisplay) {
            const commentCount = parseInt(
              commentCountDisplay.innerText.replace(/\,/g, "")
            );
            if (isNaN(commentCount)) {
              console.error("wrong comment count format");
              return;
            }
            commentCountDisplay.innerText = (commentCount - 1).toLocaleString();
          }
        }
        return;
      });
      headerEl.append(deleteIcon);
    }
    headerEl.append(avatarFrameEl);
    headerEl.append(infoEl);
    liEl.append(headerEl);
    liEl.append(textEl);
    commentsUlEl.prepend(liEl);
  }
};

export const clearCommentBox = () => {
  preparedCommentBox && (preparedCommentBox.innerHTML = "");
};

const clearComments = () => {
  commentsUlEl && (commentsUlEl.innerHTML = "");
};

const initEditorReadOnly = async (storyId: string) => {
  const { data: blocks } = await Axios.get(BASE_URL + `/api/blocks/${storyId}`);
  const header = blocks.shift();
  const headerEditor = new EditorJS({
    //   readOnly: true, // This occurs error on 2.19.0 version. It's on github issue https://github.com/codex-team/editor.js/issues/1400
    holder: "editor-readOnly__header",
    tools: {
      header: {
        class: Header,
        config: {
          levels: [2, 4, 6],
        },
      },
      code: CodeTool,
      image: {
        class: ImageTool,
      },
    },
    data: { blocks: [header] },
  });
  const bodyEditor = new EditorJS({
    //   readOnly: true, // This occurs error on 2.19.0 version. It's on github issue https://github.com/codex-team/editor.js/issues/1400
    holder: "editor-readOnly__body",
    tools: {
      header: {
        class: Header,
        config: {
          levels: [2, 4, 6],
        },
      },
      code: CodeTool,
      image: {
        class: ImageTool,
      },
    },
    data: { blocks },
  });
  headerEditor.isReady.then(async () => {
    await headerEditor.readOnly.toggle(true);
    const x = editorReadOnlyHeader?.querySelector(".codex-editor__redactor") as
      | HTMLElement
      | null
      | undefined;
    console.log(x);
    x!.style.paddingBottom = "1rem";
  });
  bodyEditor.isReady.then(async () => {
    await bodyEditor.readOnly.toggle(true);
  });
};

const closeCommentDrawer = () => {
  commentDrawer && (commentDrawer.style.right = "-26rem");
  filterBlack && (filterBlack.style.backgroundColor = "transparent");
  filterBlack && (filterBlack.style.pointerEvents = "none");
  clearCommentBox();
  clearComments();
};

const openCommentDrawer = () => {
  commentDrawer && (commentDrawer.style.right = "0rem");
  filterBlack && (filterBlack.style.backgroundColor = "black");
  filterBlack && (filterBlack.style.pointerEvents = "all");
};

const getComments = async () => {
  const splitedPath = document.location.pathname.split("read");
  const storyId = splitedPath[1].replace(/[/]/g, "");
  const { status, data } = await Axios.get(
    BASE_URL + `/api/comment/${storyId}`
  );
  if (status === 200) {
    data.forEach((comment: any) => {
      drawNewComment(comment);
    });
  }
};

const likeOrUnlike = async () => {
  const splitedPath = document.location.pathname.split("read");
  const storyId = splitedPath[1].replace(/[/]/g, "");
  const childIcon = likedContainer?.querySelector("i");
  const childSpan = likedContainer?.querySelector("span");
  if (childIcon && childSpan) {
    const likedCount = parseInt(childSpan.innerText.replace(/\,/, ""));
    if (isNaN(likedCount)) {
      console.error("wrong like count format");
      return;
    }
    let plusMinus;
    if (childIcon.className.includes("false")) {
      plusMinus = 1;
    } else if (childIcon.className.includes("true")) {
      plusMinus = -1;
    } else {
      return;
    }
    try {
      const { status } = await Axios.post(
        BASE_URL + `/api/like/${storyId}/${plusMinus}`
      );
      if (status === 200) {
        if (childIcon.className.includes("false")) {
          childSpan.innerText = (likedCount + 1).toLocaleString();
          childIcon.className = childIcon.className
            .replace("false", "true")
            .replace("far", "fas");
        } else if (childIcon.className.includes("true")) {
          childSpan.innerText = (likedCount - 1).toLocaleString();
          childIcon.className = childIcon.className
            .replace("true", "false")
            .replace("fas", "far");
        } else {
          return;
        }
      }
    } catch {}
  }
};

const handleBookmark = async () => {
  const splitedPath = document.location.pathname.split("read");
  const storyId = splitedPath[1].replace(/[/]/g, "");
  const childIcon = bookmarkContainer?.querySelector("i");
  if (childIcon) {
    if (childIcon.className.includes("false")) {
      const { status } = await Axios.post(
        BASE_URL + `/api/bookmark/${storyId}`
      );
      if (status === 200) {
        childIcon.className = childIcon.className
          .replace("false", "true")
          .replace("far", "fas");
      }
    } else if (childIcon.className.includes("true")) {
      const { status } = await Axios.delete(
        BASE_URL + `/api/bookmark/${storyId}`
      );
      if (status === 200) {
        childIcon.className = childIcon.className
          .replace("true", "false")
          .replace("fas", "far");
      }
    }
  }
};

const onFollowBtnClick = async (e: Event) => {
  const followBtn = e.currentTarget as HTMLButtonElement;
  const authorId = followBtn.closest("header")?.id;
  if (authorId) {
    const { status } = await Axios.post(BASE_URL + `/api/follow/${authorId}`);
    if (status === 200) {
      followBtn.className = followBtn.className.replace("follow", "following");
      followBtn.innerText = "Following";
      followBtn.removeEventListener("click", onFollowBtnClick);
      followBtn.addEventListener("click", onFollowingBtnClick);
    }
  }
};

const onFollowingBtnClick = async (e: Event) => {
  const followingBtn = e.currentTarget as HTMLButtonElement;
  const authorId = followingBtn.closest("header")?.id;
  if (authorId) {
    const isConfirmed = confirm("Unfollow this author?");
    if (!isConfirmed) {
      return;
    }
    const { status } = await Axios.post(BASE_URL + `/api/unfollow/${authorId}`);
    if (status === 200) {
      followingBtn.className = followingBtn.className.replace(
        "following",
        "follow"
      );
      followingBtn.innerText = "Follow";
      followingBtn.removeEventListener("click", onFollowingBtnClick);
      followingBtn.addEventListener("click", onFollowBtnClick);
    }
  }
};

const initReadStory = async () => {
  if (BASE_URL) {
    const params = document.location.pathname.split(BASE_URL)[0].split("/");
    if (params[1] === "read") {
      const storyId = params[2];
      await initEditorReadOnly(storyId);
    }
  }
  filterBlack?.addEventListener("click", closeCommentDrawer);
  seeCommentDiv?.addEventListener("click", openCommentDrawer);
  seeCommentDiv?.addEventListener("click", getComments);
  commentDrawerCloseIcon?.addEventListener("click", closeCommentDrawer);
  likedContainer?.addEventListener("click", likeOrUnlike);
  bookmarkContainer?.addEventListener("click", handleBookmark);
  followBtn?.addEventListener("click", onFollowBtnClick);
  followingBtn?.addEventListener("click", onFollowingBtnClick);
};

initReadStory();
