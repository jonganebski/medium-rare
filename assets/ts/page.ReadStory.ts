import EditorJS, { LogLevels, OutputBlockData } from "@editorjs/editorjs";
import Header from "@editorjs/header";
import CodeTool from "@editorjs/code";
import ImageTool from "@editorjs/image";
import { BASE_URL } from "./constants";
import Axios from "axios";
import {
  commentsUlEl,
  preparedCommentBox,
  editorReadOnlyHeader,
  commentDrawer,
  filterBlack,
  likedContainer,
  bookmarkContainer,
  seeCommentDiv,
  commentDrawerCloseIcon,
  followBtn,
  followingBtn,
  readTimeSpan,
  deleteStoryBtn,
} from "./elements.readStory";
import { readPageFollowBtnClick, readPageFollowingBtnClick } from "./follow";
import { deleteComment } from "./deleteComment";
import { deleteStory } from "./deleteStory";
import { getIdParam } from "./helper";

const getFomattedCommentDate = (createdAt: any): string => {
  const createdAtNum = +createdAt;
  if (!isNaN(createdAtNum)) {
    const inMilliSec = createdAtNum * 1000;
    const unixNow = new Date().getTime();
    const diff = unixNow - inMilliSec;
    const diffMin = Math.ceil(diff / 60 / 1000);
    if (diffMin === 1) {
      return `just now`;
    }
    if (diffMin < 60) {
      return `${diffMin} minutes ago`;
    }
    return new Date(inMilliSec).toLocaleString();
  } else {
    return "unknown";
  }
};

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
    const commentDate = getFomattedCommentDate(createdAt);
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
    timestampEl.innerText = commentDate;
    usernameEl.innerText = username;
    avatarImgEl.src = avatarUrl;
    infoEl.append(usernameEl);
    infoEl.append(timestampEl);
    avatarFrameEl.append(avatarImgEl);
    if (isAuthorized) {
      deleteIcon = document.createElement("i");
      deleteIcon.className = "far fa-trash-alt comment__delete-icon";
      deleteIcon.addEventListener("click", () =>
        deleteComment(liEl, commentId)
      );
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

const computeAndPasteReadTime = (blocks: OutputBlockData[]) => {
  if (!readTimeSpan) {
    return;
  }
  let wordCount = 0;
  blocks.forEach((block) => {
    if (block.type === "paragraph") {
      const words: string[] = block.data.text.split(" ");
      wordCount += words.length;
    }
    if (block.type === "code") {
      let words: string[] = block.data.code?.split(" ");
      words = words?.filter(
        (word) =>
          word !== "" &&
          !word.includes("=") &&
          word != "()" &&
          word != "(" &&
          word != ")" &&
          word != "{" &&
          word != "}" &&
          word != "<" &&
          word != ">"
      );
      wordCount += words?.length;
    }
  });
  const readTimeMinute = Math.ceil(wordCount / 200);
  !isNaN(readTimeMinute) &&
    (readTimeSpan.innerText = `${readTimeMinute} min read`);
};

const overrideEditorJsStyle = () => {
  const x = editorReadOnlyHeader?.querySelector(".codex-editor__redactor") as
    | HTMLElement
    | null
    | undefined;
  x!.style.paddingBottom = "1rem";
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
    logLevel: LogLevels?.ERROR ?? "ERROR",
  });
  headerEditor.isReady.then(async () => {
    await headerEditor.readOnly.toggle(true);
    overrideEditorJsStyle();
  });
  bodyEditor.isReady.then(async () => {
    await bodyEditor.readOnly.toggle(true);
    computeAndPasteReadTime(blocks);
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
  const storyId = getIdParam("read-story");
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
  const storyId = getIdParam("read-story");
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
      const { status } = await Axios.patch(
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
  const storyId = getIdParam("read-story");
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
      const { status } = await Axios.patch(
        BASE_URL + `/api/disbookmark/${storyId}`
      );
      if (status === 200) {
        childIcon.className = childIcon.className
          .replace("true", "false")
          .replace("fas", "far");
      }
    }
  }
};

const initReadStory = async () => {
  if (BASE_URL) {
    const params = document.location.pathname.split(BASE_URL)[0].split("/");
    if (params[1] === "read-story") {
      const storyId = getIdParam("read-story");
      await initEditorReadOnly(storyId);
    }
  }
  filterBlack?.addEventListener("click", closeCommentDrawer);
  seeCommentDiv?.addEventListener("click", openCommentDrawer);
  seeCommentDiv?.addEventListener("click", getComments);
  commentDrawerCloseIcon?.addEventListener("click", closeCommentDrawer);
  likedContainer?.addEventListener("click", likeOrUnlike);
  bookmarkContainer?.addEventListener("click", handleBookmark);
  followBtn?.addEventListener("click", readPageFollowBtnClick);
  followingBtn?.addEventListener("click", readPageFollowingBtnClick);
  deleteStoryBtn?.addEventListener("click", deleteStory);
};

initReadStory();
