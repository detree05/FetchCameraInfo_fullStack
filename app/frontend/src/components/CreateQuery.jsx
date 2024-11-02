import { verbose, cameraName, projectName, username, password } from './DataHandler';

function CreateQuery() {
    
    return {
        Verbose: verbose,
        CameraName: cameraName,
        ProjectName: projectName,
        Username: username,
        Password: password
    }
  }
    
  export default CreateQuery