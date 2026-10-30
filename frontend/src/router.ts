import { createBrowserRouter, type RouteObject } from "react-router";

import HikeLog from "./pages/HikeLog/HikeLog/HikeLog";
import Login from "./pages/Login/Login";

const routes: RouteObject[] = [
  {
    Component: HikeLog,
    path: "/",
  },
  {
    Component: Login,
    path: "/login",
  },
];
const router = createBrowserRouter(routes);

export default router;
