cmake_minimum_required(VERSION 2.6 FATAL_ERROR)

project(microservice NONE)

## Install ####################################################################
###############################################################################

if(NOT DEFINED TARANTOOL_INSTALL_LUADIR)
    set(TARANTOOL_INSTALL_LUADIR "${PROJECT_SOURCE_DIR}/.rocks/share/tarantool")
endif()

install(
  FILES ${CMAKE_CURRENT_SOURCE_DIR}/init.lua
  DESTINATION ${TARANTOOL_INSTALL_LUADIR}/${PROJECT_NAME}
)
