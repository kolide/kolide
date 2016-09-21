import React, { PropTypes } from 'react';
import radium from 'radium';
import componentStyles from './styles';

const Avatar = ({ size, style, user }) => {
  const { gravatarURL } = user;

  return (
    <img
      alt="User avatar"
      src={gravatarURL}
      style={[componentStyles(size), style]}
    />
  );
};

Avatar.propTypes = {
  size: PropTypes.string,
  style: PropTypes.object,
  user: PropTypes.object.isRequired,
};

export default radium(Avatar);
