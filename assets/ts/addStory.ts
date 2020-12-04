// import Axios from "axios";
import marked from "marked";

// const submitBtn = document.getElementById("test-submit");
const resultDiv = document.querySelector("._result");
const html = marked("# Marked in Node.js\n\nRendered by **marked**.");
const addStory = () => {
  resultDiv!.innerHTML = html;
};

addStory();
