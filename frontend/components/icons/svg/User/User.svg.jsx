import React, { Component } from 'react';
import base, { basePropTypes } from '../base';

export class User extends Component {
  static propTypes = {
    ...basePropTypes,
  };

  static variants = {
    default: (
      <g stroke="none" strokeWidth="1" fill="none" fillRule="evenodd">
        <g transform="translate(-672.000000, -519.000000)" fill="#EAEEFB" stroke="#B9C2E4">
          <g transform="translate(272.000000, 340.000000)">
            <g transform="translate(41.000000, 170.000000)">
              <g transform="translate(360.000000, 9.500000)">
                <path d="M5,12 C7.76142375,12 10,10.4198285 10,8.47058824 C10,6.52134794 7.76142375,4.94117647 5,4.94117647 C2.23857625,4.94117647 0,6.52134794 0,8.47058824 C0,10.4198285 2.23857625,12 5,12 Z" />
                <ellipse cx="5" cy="3.52941176" rx="3.57142857" ry="3.52941176" />
              </g>
            </g>
          </g>
        </g>
      </g>
    ),
    colored: (
      <g stroke="none" strokeWidth="1" fill="none" fillRule="evenodd">
        <g transform="translate(-672.000000, -519.000000)" fill="#EED6FF" stroke="#C48DED">
          <g transform="translate(272.000000, 340.000000)">
            <g transform="translate(41.000000, 149.000000)">
              <g transform="translate(360.000000, 30.500000)">
                <path d="M5,12 C7.76142375,12 10,10.4198285 10,8.47058824 C10,6.52134794 7.76142375,4.94117647 5,4.94117647 C2.23857625,4.94117647 0,6.52134794 0,8.47058824 C0,10.4198285 2.23857625,12 5,12 Z" />
                <ellipse cx="5" cy="3.52941176" rx="3.57142857" ry="3.52941176" />
              </g>
            </g>
          </g>
        </g>
      </g>
    ),
  };

  render () {
    const { alt, style, variant } = this.props;

    if (variant === 'circle') {
      return (
        <svg
          width="98px"
          height="99px"
          viewBox="0 0 98 99"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          xmlnsXlink="http://www.w3.org/1999/xlink"
          alt={alt}
          style={style}
        >
          {User.variants[variant]}
        </svg>
      );
    }

    return (
      <svg
        width="12px"
        height="13px"
        viewBox="0 0 12 13"
        version="1.1"
        xmlns="http://www.w3.org/2000/svg"
        xmlnsXlink="http://www.w3.org/1999/xlink"
        alt={alt}
        style={style}
      >
        {User.variants[variant]}
      </svg>
    );
  }
}

export default base(User);
