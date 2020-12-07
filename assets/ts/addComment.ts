import Axios from "axios";
import { BASE_URL } from "./constants";
import {
  commentDrawer,
  preparedCommentBox,
  commentCountDisplay,
} from "./elements.readStory";
import { drawNewComment, clearCommentBox } from "./page.ReadStory";

const cancelCommentBtn = commentDrawer?.querySelector(
  ".add-comment__cancel-btn"
);
const addCommentBtn = commentDrawer?.querySelector(".add-comment__add-btn");

const initAddComment = () => {
  addCommentBtn?.addEventListener("click", async () => {
    if (preparedCommentBox && commentCountDisplay) {
      const commentText = preparedCommentBox.innerText;
      const commentCount = parseInt(
        commentCountDisplay.innerText.replace(/\,/g, "")
      );
      if (isNaN(commentCount)) {
        console.error("wrong comment count format");
        return;
      }
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
          commentCountDisplay.innerText = (commentCount + 1).toLocaleString();
          drawNewComment(comment);
        }
      } catch {}
    }
  });
  cancelCommentBtn?.addEventListener("click", clearCommentBox);
};

initAddComment();
