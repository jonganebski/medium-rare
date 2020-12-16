import Axios from "axios";
import { BASE_URL } from "./constants";

export const deleteStory = async () => {
  const isConfirmed = confirm(
    `You are removing this story permanently.
Are you sure?`
  );
  if (!isConfirmed) {
    return;
  }
  try {
    const splitedPath = document.location.pathname.split("read-story");
    const storyId = splitedPath[1].replace(/[/]/g, "");
    const { status } = await Axios.delete(`/api/story/${storyId}`);
    if (status < 300) {
      document.location.href = BASE_URL + "/me/stories";
    }
  } catch {
    alert("Failed to delete story. Please try again.");
    document.location.href = BASE_URL + "/me/stories";
  }
};
