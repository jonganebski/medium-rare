import Axios from "axios";
import { followingPageHeader } from "./elements.followings";
import { followersCountDisplay } from "./elements.readStory";

export const readPageFollowBtnClick = async (e: Event) => {
  const followBtn = e.currentTarget as HTMLButtonElement;
  const authorId = followBtn.closest("header")?.id;
  if (authorId) {
    try {
      const { status } = await Axios.patch(`/api/follow/${authorId}`);
      if (status < 300) {
        const followersCount = followersCountDisplay?.innerText
          .replace("Followers", "")
          .replace(/[,]/g, "")
          .trim();
        const followersCountLink = followersCountDisplay?.querySelector("a");
        if (followersCountLink && followersCount && !isNaN(+followersCount)) {
          followersCountLink.innerText =
            (+followersCount + 1).toLocaleString() + " Followers";
        }
        followBtn.className = followBtn.className.replace(
          "follow",
          "following"
        );
        followBtn.innerText = "Following";
        followBtn.removeEventListener("click", readPageFollowBtnClick);
        followBtn.addEventListener("click", readPageFollowingBtnClick);
      }
    } catch {
      alert("Failed to follow. Please try again.");
    }
  }
};

export const readPageFollowingBtnClick = async (e: Event) => {
  const followingBtn = e.currentTarget as HTMLButtonElement;
  const authorId = followingBtn.closest("header")?.id;
  if (authorId) {
    const isConfirmed = confirm("Unfollow this author?");
    if (!isConfirmed) {
      return;
    }
    try {
      const { status } = await Axios.patch(`/api/unfollow/${authorId}`);
      if (status < 300) {
        const followersCountLink = followersCountDisplay?.querySelector("a");
        const followersCount = followersCountLink?.innerText
          .replace("Followers", "")
          .replace(/[,]/g, "")
          .trim();
        if (
          followersCountLink &&
          followersCount &&
          !isNaN(+followersCount) &&
          +followersCount !== 0
        ) {
          followersCountLink.innerText =
            (+followersCount - 1).toLocaleString() + " Followers";
        }
        followingBtn.className = followingBtn.className.replace(
          "following",
          "follow"
        );
        followingBtn.innerText = "Follow";
        followingBtn.removeEventListener("click", readPageFollowingBtnClick);
        followingBtn.addEventListener("click", readPageFollowBtnClick);
      }
    } catch {
      alert("Failed to unfollow. Please try again.");
    }
  }
};

export const followersPageFollowBtnClick = async (e: Event) => {
  const followBtn = e.currentTarget as HTMLButtonElement | null;
  const authorId = followBtn?.closest("li")?.id;
  if (!authorId || !followBtn) {
    return;
  }
  try {
    const { status } = await Axios.patch(`/api/follow/${authorId}`);
    if (status < 300) {
      followBtn.className = "userCard__following-btn";
      followBtn.innerText = "Following";
      followBtn?.removeEventListener("click", followersPageFollowBtnClick);
      followBtn?.addEventListener("click", followersPageFollowingBtnClick);
    }
  } catch {
    alert("Failed to follow. Please try again.");
  }
};

export const followersPageFollowingBtnClick = async (e: Event) => {
  const followingBtn = e.currentTarget as HTMLButtonElement | null;
  const authorId = followingBtn?.closest("li")?.id;
  if (!authorId || !followingBtn) {
    return;
  }
  const isConfirmed = confirm("Unfollow this author?");
  if (!isConfirmed) {
    return;
  }
  try {
    const { status } = await Axios.patch(`/api/unfollow/${authorId}`);
    if (status === 200) {
      followingBtn.className = "userCard__follow-btn";
      followingBtn.innerText = "Follow";
      followingBtn.removeEventListener("click", followersPageFollowingBtnClick);
      followingBtn.addEventListener("click", followersPageFollowBtnClick);
    }
  } catch {
    alert("Failed to unfollow. Please try again.");
  }
};

export const followingsPageUnfollowBtnClick = async (e: Event) => {
  const unfollowBtn = e.currentTarget as HTMLButtonElement | null;
  const userCard = unfollowBtn?.closest("li");
  const authorId = userCard?.id;
  if (!authorId || !userCard || !unfollowBtn || !followingPageHeader) {
    return;
  }
  const isConfirmed = confirm("Unfollow this author?");
  if (!isConfirmed) {
    return;
  }
  try {
    const { status } = await Axios.patch(`/api/unfollow/${authorId}`);
    if (status < 300) {
      userCard.remove();
      const prevCount = followingPageHeader.innerText
        .replace("You are following", "")
        .replace("medium-rares.", "")
        .replace(/[,]/g, "")
        .trim();
      console.log(prevCount);
      if (!prevCount || isNaN(+prevCount)) {
        return;
      }
      followingPageHeader.innerText =
        "You are following " +
        (+prevCount - 1).toLocaleString() +
        " medium-rares.";
    }
  } catch {
    alert("Failed to follow. Please try again.");
  }
};

export const userBioFollowBtnClick = async (e: Event) => {
  const target = e.target as HTMLButtonElement | null;
  const parentEl = target?.parentElement;
  const targetUserId = parentEl?.id;
  if (!targetUserId) {
    return;
  }
  try {
    const { status } = await Axios.patch(`/api/follow/${targetUserId}`);
    if (status === 200) {
      target && (target.className = "userBio__following-btn");
      target && (target.innerText = "Following");
      target && target.removeEventListener("click", userBioFollowBtnClick);
      target && target.addEventListener("click", userBioFollowingBtnClick);
    }
  } catch {
    alert("Failed to follow. Please try again.");
  }
};

export const userBioFollowingBtnClick = async (e: Event) => {
  const target = e.target as HTMLButtonElement | null;
  const parentEl = target?.parentElement;
  const targetUserId = parentEl?.id;
  if (!targetUserId) {
    return;
  }
  const isConfirmed = confirm("Unfollow this user?");
  if (!isConfirmed) {
    return;
  }
  try {
    const { status } = await Axios.patch(`/api/unfollow/${targetUserId}`);
    if (status < 300) {
      target && (target.className = "userBio__follow-btn");
      target && (target.innerText = "Follow");
      target && target.removeEventListener("click", userBioFollowingBtnClick);
      target && target.addEventListener("click", userBioFollowBtnClick);
    }
  } catch {
    alert("Failed to unfollow. Please try again.");
  }
};
