import { userBioFollowBtnClick, userBioFollowingBtnClick } from "./follow";

const userBio = document.getElementById("userBio");
const userBioFollowBtn = userBio?.querySelector(".userBio__follow-btn");
const userBioFollowingBtn = document.querySelector(".userBio__following-btn");

const init = () => {
  userBioFollowBtn?.addEventListener("click", userBioFollowBtnClick);
  userBioFollowingBtn?.addEventListener("click", userBioFollowingBtnClick);
};

init();
