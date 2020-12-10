import Axios from "axios";
import { BASE_URL } from "./constants";
import {
  addCommentBtn,
  cancelCommentBtn,
  commentCountDisplay,
  commentDrawerCommentCount,
  preparedCommentBox,
} from "./elements.readStory";
import { clearCommentBox, drawNewComment } from "./page.ReadStory";

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
          if (commentDrawerCommentCount) {
            const prevCount = +commentDrawerCommentCount.innerText;
            if (!isNaN(prevCount)) {
              commentDrawerCommentCount.innerText = prevCount + 1 + "";
            }
          }
          preparedCommentBox.innerText = "";
          commentCountDisplay.innerText = (commentCount + 1).toLocaleString();
          drawNewComment(comment);
        }
      } catch {}
    }
  });
  cancelCommentBtn?.addEventListener("click", clearCommentBox);
};

initAddComment();
