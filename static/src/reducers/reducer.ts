import {HYDRATE} from 'next-redux-wrapper'
import {AnyAction, combineReducers, Reducer} from 'redux';
import {counter, counterInitialState} from './counter';
import {global, globalInitialState} from './global';
import {indexInitialState, indexPage} from "./index";

const combinedReducer = combineReducers({
  counter,
  global,
  indexPage,
});

export const reducer: Reducer<State, AnyAction> = (state, action) => {
  if (action.type === HYDRATE) {
    return {
      ...state, // use previous state
      ...action.payload, // apply delta from hydration
    }
  } else {
    return combinedReducer(state, action)
  }
}

export const initialState = {
  counter: counterInitialState,
  global: globalInitialState,
  indexPage: indexInitialState,
};

export type State = typeof initialState;