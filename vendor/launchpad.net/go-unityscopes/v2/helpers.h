#ifndef UNITYSCOPE_HELPERS_H
#define UNITYSCOPE_HELPERS_H

#include <string>
#include <vector>

typedef struct StrData StrData;

namespace gounityscopes {
namespace internal {

std::string from_gostring(const StrData str);
std::vector<std::string> split_strings(const StrData str);
void *as_bytes(const std::string &str, int *length);

}
}

#endif
