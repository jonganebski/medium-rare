import Axios from "axios";
import {
  commentCountDisplay,
  commentDrawerCommentCount,
} from "./elements.readStory";

export const deleteComment = async (
  commentContainer: HTMLLIElement,
  commentId: string
) => {
  try {
    const { status } = await Axios.delete(`/api/comment/${commentId}`);
    if (status < 300) {
      commentContainer.remove();
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
      if (commentDrawerCommentCount) {
        const prevCount = +commentDrawerCommentCount.innerText;
        if (!isNaN(prevCount) && prevCount !== 0) {
          commentDrawerCommentCount.innerText = prevCount - 1 + "";
        }
      }
    }
  } catch {
    alert("Failed to delete comment. Please try again.");
  }
};
