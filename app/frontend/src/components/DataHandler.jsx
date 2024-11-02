var verbose = false;
var projectName = "";
var cameraName = "";
var username = "";
var password = "";

const setVerbose = (verboseParam) => {
    verbose = verboseParam
};
const setProjectName = (projectNameParam) => {
    projectName = projectNameParam
};
const setCameraName = (cameraNameParam) => {
    cameraName = cameraNameParam
};
const setUsername = (usernameParam) => {
    username = usernameParam
};
const setPassword = (passwordParam) => {
    password = passwordParam
};

export { setVerbose, setProjectName, setCameraName, setUsername, setPassword, verbose, projectName, cameraName, username, password }