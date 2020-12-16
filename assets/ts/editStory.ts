import Axios from "axios";
import { BASE_URL } from "./constants";
import { getIdParam } from "./helper";
import { useEditor } from "./useEditor";

const init = async () => {
  const storyId = getIdParam("edit-story");
  try {
    const { data: blocks } = await Axios.get(`/api/blocks/${storyId}`);
    useEditor("editor-edit", "", blocks);
  } catch {
    alert("Failed to load story. Please try again.");
    document.location.href = BASE_URL + `/read-story/${storyId}`;
  }
};

if (document.location.pathname.includes("edit-story")) {
  init();
}
