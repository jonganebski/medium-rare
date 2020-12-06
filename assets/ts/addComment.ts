import Axios, { AxiosResponse } from "axios";
import { BASE_URL } from "./constants";

const commentDrawer = document.getElementById("drawer-comment");
const preparedCommentBox = commentDrawer?.querySelector(
  ".add-comment__text"
) as HTMLParagraphElement | null;
const cancelCommentBtn = commentDrawer?.querySelector(
  ".add-comment__cancel-btn"
);
const addCommentBtn = commentDrawer?.querySelector(".add-comment__add-btn");
const commentsUlEl = commentDrawer?.querySelector("ul");

const drawNewComment = (response: AxiosResponse<any>) => {
  if (commentsUlEl) {
    const { avatarUrl, createdAt, text, username } = response.data;
    const textEl = document.createElement("p");
    const timestampEl = document.createElement("div");
    const usernameEl = document.createElement("div");
    const infoEl = document.createElement("div");
    const avatarImgEl = document.createElement("img");
    const avatarFrameEl = document.createElement("div");
    const deleteIcon = document.createElement("i");
    const headerEl = document.createElement("header");
    const liEl = document.createElement("li");
    textEl.className = "comment__text";
    timestampEl.className = "comment__timestamp";
    usernameEl.className = "comment__creator-name";
    infoEl.className = "comment__info _flex-c-sb";
    avatarImgEl.className = "_avatar-img";
    avatarFrameEl.className = "_avatar-frame";
    deleteIcon.className = "far fa-trash-alt comment__delete-icon";
    headerEl.className = "comment__header _flex-cs";
    liEl.className = "drawer-comment__comment-container";
    textEl.innerText = text;
    timestampEl.innerText = createdAt;
    usernameEl.innerText = username;
    avatarImgEl.src = avatarUrl;
    infoEl.append(usernameEl);
    infoEl.append(timestampEl);
    avatarFrameEl.append(avatarImgEl);
    headerEl.append(deleteIcon);
    headerEl.append(avatarFrameEl);
    headerEl.append(infoEl);
    liEl.append(headerEl);
    liEl.append(textEl);
    commentsUlEl.prepend(liEl);
  }
};

const initAddComment = () => {
  addCommentBtn?.addEventListener("click", async () => {
    if (preparedCommentBox) {
      const splitedPath = document.location.pathname.split("/");
      const storyId = splitedPath[splitedPath.length - 1];
      try {
        const response = await Axios.post(
          BASE_URL + `/api/comment/${storyId}`,
          {
            text: preparedCommentBox.innerText,
          }
        );
        if (response.status === 201) {
          drawNewComment(response);
        }
      } catch {}
    }
  });
};

initAddComment();
