import React from 'react';
import { css } from '@emotion/core';
import { Link } from '@reach/router';

const ActiveLink = ({ ...props }) => {
  return (
    <Link
      getProps={({ isCurrent }) => (isCurrent ? { className: 'current' } : {})}
      {...props}
    />
  );
};

export default function Header() {
  return (
    <div
      css={css`
        margin-top: 30px;
        @media only screen and (max-width: 1000px) {
          margin-top: 10px;
        }
        display: flex;
      `}
    >
      <div
        css={css`
          flex: 1;
          display: flex;
          border-bottom: 1px solid #ccc;
          padding: 10px 0;
          justify-content: space-between;
          align-items: flex-end;
        `}
      >
        <div
          css={css`
            font-size: 40px;
            @media only screen and (max-width: 1000px) {
              font-size: 20px;
            }
            color: #333;
            font-weight: bold;
          `}
        >
          Covid Analytics
        </div>
        <div
          css={css`
            text-transform: uppercase;
            font-size: 20px;
            display; flex;
            flex-direction: row;
            a {
              margin-left: 20px;
              @media only screen and (max-width: 1000px) {
                margin-left: 5px;
                font-size: 10px;
              }
              text-decoration: none;
              color: #1e88e5;
              padding: 4px 8px;
              &:hover {
                background-color: #eeeeee;
              }
              &.current {
                color: #333;
                cursor: default;
                pointer-events: none;
                font-weight: bold;
                &:hover {
                  background-color: white;
                }
              }
            }
          `}
        >
          <ActiveLink to="/world">World</ActiveLink>
          <ActiveLink to="/us">Us</ActiveLink>
        </div>
      </div>
    </div>
  );
}
