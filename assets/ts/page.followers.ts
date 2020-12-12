import {
  followersPagefollowBtns,
  followersPagefollowingBtns,
} from "./elements.followers";
import {
  followersPageFollowBtnClick,
  followersPageFollowingBtnClick,
} from "./follow";

const initPageFollowers = () => {
  followersPagefollowBtns?.forEach((btn) => {
    btn.addEventListener("click", followersPageFollowBtnClick);
  });
  followersPagefollowingBtns?.forEach((btn) => {
    btn.addEventListener("click", followersPageFollowingBtnClick);
  });
};

initPageFollowers();
