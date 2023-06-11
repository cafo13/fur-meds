import firebaseConfigJson from '../../firebase.config.json';
import packageJson from '../../package.json';

export const environment = {
  production: true,
  apiEndpoint: 'https://furmends-api-fh65cjqo3q-ez.a.run.app/api/v1',
  firebaseConfig: firebaseConfigJson.firebaseConfig,
  version: packageJson.version,
};
