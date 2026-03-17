const e=async({cookies:e})=>({isAuthenticated:!!e.get("access_token")});export{e as load};
