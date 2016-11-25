import labels from './label';
import experiments from './experiments';

const namespaceInitialState = {
  name: '',
  labels: [],
  experiments: [],
  isDirty: false,
  isNew: false,
  delete: false,
  publish: false,
};

const namespace = (state = namespaceInitialState, action) => {
  switch (action.type) {
  case 'ADD_NAMESPACE':
      return { ...state, name: action.name, isDirty: true, isNew: true };
  case 'NAMESPACE_NAME':
    return { ...state, name: action.name, isDirty: true };
  case 'ADD_LABEL':
  case 'TOGGLE_LABEL':
    return { ...state, labels: labels(state.labels, action), isDirty: true };
  case 'NAMESPACE_DELETE':
    return { ...state, delete: true, isDirty: true };
  case 'TOGGLE_PUBLISH':
    return { ...state, publish: !state.publish };
  case 'ADD_EXPERIMENT':
  case 'EXPERIMENT_DELETE':
  case 'EXPERIMENT_NAME':
  case 'EXPERIMENT_NUM_SEGMENTS':
  case 'PARAM_NAME':
  case 'ADD_PARAM':
  case 'PARAM_DELETE':
  case 'TOGGLE_WEIGHTED':
  case 'ADD_CHOICE':
  case 'CHOICE_DELETE':
  case 'ADD_WEIGHT':
  case 'CLEAR_CHOICES':
    return { ...state, experiments: experiments(state.experiments, action), isDirty: true };
  default: 
    return state;
  }
}

const namespaces = (state = [], action) => {
  switch (action.type) {
  case 'NAMESPACES_LOADED':
    return action.namespaces.map(n => Object.assign({}, namespace(undefined, action), n));
  case 'ADD_NAMESPACE':
    return [...state, namespace(undefined, action)];
  case 'NAMESPACE_DELETE':
  case 'NAMESPACE_NAME':
  case 'ADD_LABEL':
  case 'TOGGLE_LABEL':
  case 'TOGGLE_PUBLISH':
  case 'ADD_EXPERIMENT':
  case 'EXPERIMENT_DELETE':
  case 'EXPERIMENT_NAME':
  case 'EXPERIMENT_NUM_SEGMENTS':
  case 'PARAM_NAME':
  case 'ADD_PARAM':
  case 'PARAM_DELETE':
  case 'TOGGLE_WEIGHTED':
  case 'ADD_CHOICE':
  case 'CHOICE_DELETE':
  case 'ADD_WEIGHT':
  case 'CLEAR_CHOICES':
    const ns = state.map(n => {
      if (n.name !== action.namespace) {
        return n;
      }
      return namespace(n, action);
    });
    return ns;
  default:
    return state;
  }
}

export default namespaces;
