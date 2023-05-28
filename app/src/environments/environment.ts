import firebaseConfigJson from '../../firebase.config.json';
import packageJson from '../../package.json';

export const environment = {
  production: false,
  // apiEndpoint: 'http://localhost:8080/api/v1',
  apiEndpoint: 'https://pets-api-fh65cjqo3q-ez.a.run.app/api/v1',
  firebaseConfig: firebaseConfigJson.firebaseConfig,
  version: packageJson.version,
};
