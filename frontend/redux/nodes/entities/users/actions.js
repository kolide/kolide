import Kolide from 'kolide';

import config from './config';

export const REQUIRE_PASSWORD_RESET_REQUEST = 'REQUIRE_PASSWORD_RESET_REQUEST';
export const REQUIRE_PASSWORD_RESET_SUCCESS = 'REQUIRE_PASSWORD_RESET_SUCCESS';
export const REQUIRE_PASSWORD_RESET_FAILURE = 'REQUIRE_PASSWORD_RESET_FAILURE';

export const requirePasswordResetRequest = () => {
  return { type: REQUIRE_PASSWORD_RESET_REQUEST };
};

export const requirePasswordResetSuccess = (user) => {
  return {
    type: REQUIRE_PASSWORD_RESET_SUCCESS,
    payload: { user },
  };
};

export const requirePasswordResetFailure = (errors) => {
  return {
    type: REQUIRE_PASSWORD_RESET_FAILURE,
    payload: { errors },
  };
};

export const requirePasswordReset = (user, { require }) => {
  return (dispatch) => {
    dispatch(requirePasswordResetRequest());

    return Kolide.requirePasswordReset(user, { require })
      .then((updatedUser) => {
        dispatch(requirePasswordResetSuccess(updatedUser));

        return updatedUser;
      })
      .catch((response) => {
        const { errors } = response;

        dispatch(requirePasswordResetFailure(errors));
        throw response;
      });
  };
};

export default { ...config.actions, requirePasswordReset };
