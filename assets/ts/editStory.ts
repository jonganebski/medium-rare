import Axios from "axios";
import { BASE_URL } from "./constants";
import { useEditor } from "./useEditor";

const editStory = async () => {
  if (BASE_URL) {
    const params = document.location.pathname.split(BASE_URL)[0].split("/");
    if (params[1] === "edit-story") {
      const storyId = params[2];
      const { data: blocks } = await Axios.get(
        BASE_URL + `/api/blocks/${storyId}`
      );
      useEditor("editor-edit", "", blocks);
    }
  }
};

editStory();
