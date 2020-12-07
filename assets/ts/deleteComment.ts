import Axios from "axios";
import { BASE_URL } from "./constants";
import { commentCountDisplay } from "./elements.readStory";

export const deleteComment = async (
  commentContainer: HTMLLIElement,
  commentId: string
) => {
  const { status } = await Axios.delete(BASE_URL + `/api/comment/${commentId}`);
  if (status === 200) {
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
  }
  return;
};
