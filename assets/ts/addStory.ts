import { INITIAL_BLOCKS } from "./constants";
import { useEditor } from "./useEditor";

const init = () => {
  useEditor("editor-write", "Write your story", INITIAL_BLOCKS);
};

if (document.location.pathname.includes("new-story")) {
  init();
}
