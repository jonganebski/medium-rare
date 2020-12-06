import Axios from "axios";
import { BASE_URL } from "./constants";
import {
  clearCommentBox,
  commentDrawer,
  drawNewComment,
  preparedCommentBox,
} from "./readStory";

const cancelCommentBtn = commentDrawer?.querySelector(
  ".add-comment__cancel-btn"
);
const addCommentBtn = commentDrawer?.querySelector(".add-comment__add-btn");

const initAddComment = () => {
  addCommentBtn?.addEventListener("click", async () => {
    if (preparedCommentBox) {
      const commentText = preparedCommentBox.innerText;
      if (commentText.length === 0) {
        return;
      }
      const splitedPath = document.location.pathname.split("read");
      const storyId = splitedPath[1].replace(/[/]/g, "");
      try {
        const { status, data: comment } = await Axios.post(
          BASE_URL + `/api/comment/${storyId}`,
          {
            text: commentText,
          }
        );
        if (status === 201) {
          drawNewComment(comment);
        }
      } catch {}
    }
  });
  cancelCommentBtn?.addEventListener("click", clearCommentBox);
};

initAddComment();
