import Axios from "axios";
import { BASE_URL } from "./constants";

export const deleteStory = async () => {
  const isConfirmed = confirm(
    `You are removing this story permanently. 
    Are you sure?`
  );
  if (isConfirmed) {
    const splitedPath = document.location.pathname.split("read");
    const storyId = splitedPath[1].replace(/[/]/g, "");
    const { status } = await Axios.delete(BASE_URL + `/api/story/${storyId}`);
    if (status === 204) {
      document.location.href = BASE_URL + "/me/stories";
    }
  }
};
