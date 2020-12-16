import Axios from "axios";
import {
  addCommentBtn,
  cancelCommentBtn,
  commentCountDisplay,
  commentDrawerCommentCount,
  preparedCommentBox,
} from "./elements.readStory";
import { getIdParam } from "./helper";
import { clearCommentBox, drawNewComment } from "./page.ReadStory";

const addComment = async () => {
  if (preparedCommentBox && commentCountDisplay && commentDrawerCommentCount) {
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
    const storyId = getIdParam("read-story");
    try {
      const { status, data: comment } = await Axios.post(
        `/api/comment/${storyId}`,
        {
          text: commentText,
        }
      );
      if (status < 300) {
        const prevCount = commentDrawerCommentCount.innerText;
        if (prevCount && !isNaN(+prevCount)) {
          commentDrawerCommentCount.innerText = +prevCount + 1 + "";
        }
        preparedCommentBox.innerText = "";
        commentCountDisplay.innerText = (commentCount + 1).toLocaleString();
        drawNewComment(comment);
      }
    } catch {
      alert("Failed to add comment. Please try again.");
    }
  }
};

const init = () => {
  addCommentBtn?.addEventListener("click", addComment);
  cancelCommentBtn?.addEventListener("click", clearCommentBox);
};

init();
