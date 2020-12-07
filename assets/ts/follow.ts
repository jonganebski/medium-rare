import Axios from "axios";
import { BASE_URL } from "./constants";

export const onFollowBtnClick = async (e: Event) => {
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

export const onFollowingBtnClick = async (e: Event) => {
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
